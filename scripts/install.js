#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const REPO = 'cloudboy-jh/annotr';
const BINARY = 'annotr';

function getPlatform() {
  const platform = process.platform;
  switch (platform) {
    case 'darwin': return 'darwin';
    case 'linux': return 'linux';
    case 'win32': return 'windows';
    default: throw new Error(`Unsupported platform: ${platform}`);
  }
}

function getArch() {
  const arch = process.arch;
  switch (arch) {
    case 'x64': return 'amd64';
    case 'arm64': return 'arm64';
    default: throw new Error(`Unsupported architecture: ${arch}`);
  }
}

function getLatestRelease() {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'api.github.com',
      path: `/repos/${REPO}/releases/latest`,
      headers: { 'User-Agent': 'annotr-npm-installer' }
    };

    https.get(options, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          const json = JSON.parse(data);
          resolve(json.tag_name || 'v0.1.0');
        } catch {
          resolve('v0.1.0');
        }
      });
    }).on('error', () => resolve('v0.1.0'));
  });
}

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    
    const request = (url) => {
      https.get(url, (res) => {
        if (res.statusCode === 302 || res.statusCode === 301) {
          request(res.headers.location);
          return;
        }
        
        if (res.statusCode !== 200) {
          reject(new Error(`Download failed: ${res.statusCode}`));
          return;
        }

        res.pipe(file);
        file.on('finish', () => {
          file.close();
          resolve();
        });
      }).on('error', reject);
    };

    request(url);
  });
}

async function main() {
  try {
    const platform = getPlatform();
    const arch = getArch();
    const version = await getLatestRelease();
    
    const ext = platform === 'windows' ? '.exe' : '';
    const filename = `${BINARY}-${platform}-${arch}${ext}`;
    const url = `https://github.com/${REPO}/releases/download/${version}/${filename}`;
    
    const binDir = path.join(__dirname, '..', 'bin');
    const binPath = path.join(binDir, BINARY + ext);

    console.log(`Downloading annotr ${version} for ${platform}-${arch}...`);
    
    await downloadFile(url, binPath);
    
    if (platform !== 'windows') {
      fs.chmodSync(binPath, 0o755);
    }

    console.log('annotr installed successfully!');
    console.log('Run "annotr init" to get started.');
  } catch (err) {
    console.error('Failed to install annotr:', err.message);
    console.error('');
    console.error('You can install manually:');
    console.error('  curl -fsSL https://raw.githubusercontent.com/cloudboy-jh/annotr/main/install.sh | sh');
    process.exit(1);
  }
}

main();
