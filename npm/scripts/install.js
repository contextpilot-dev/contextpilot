#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const zlib = require('zlib');
const tar = require('tar');

const REPO = 'jitin-nhz/contextpilot';
const BINARY_NAME = 'contextpilot';

// Detect platform
function getPlatform() {
  const platform = process.platform;
  switch (platform) {
    case 'darwin': return 'darwin';
    case 'linux': return 'linux';
    case 'win32': return 'windows';
    default: throw new Error(`Unsupported platform: ${platform}`);
  }
}

// Detect architecture
function getArch() {
  const arch = process.arch;
  switch (arch) {
    case 'x64': return 'amd64';
    case 'arm64': return 'arm64';
    default: throw new Error(`Unsupported architecture: ${arch}`);
  }
}

// Get latest release version from GitHub
async function getLatestVersion() {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'api.github.com',
      path: `/repos/${REPO}/releases/latest`,
      headers: { 'User-Agent': 'contextpilot-npm-installer' }
    };

    https.get(options, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          const release = JSON.parse(data);
          resolve(release.tag_name);
        } catch (e) {
          reject(new Error('Failed to parse release info'));
        }
      });
    }).on('error', reject);
  });
}

// Download file
function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    
    const request = (url) => {
      https.get(url, (res) => {
        // Handle redirects
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
      }).on('error', (err) => {
        fs.unlink(dest, () => {});
        reject(err);
      });
    };
    
    request(url);
  });
}

// Extract tar.gz
async function extractTarGz(file, dest) {
  return new Promise((resolve, reject) => {
    fs.createReadStream(file)
      .pipe(zlib.createGunzip())
      .pipe(tar.extract({ cwd: dest }))
      .on('finish', resolve)
      .on('error', reject);
  });
}

// Main install function
async function install() {
  console.log('📦 Installing ContextPilot...');
  
  const platform = getPlatform();
  const arch = getArch();
  
  console.log(`   Platform: ${platform}/${arch}`);
  
  try {
    const version = await getLatestVersion();
    console.log(`   Version: ${version}`);
    
    const filename = `${BINARY_NAME}-${platform}-${arch}${platform === 'windows' ? '.zip' : '.tar.gz'}`;
    const downloadUrl = `https://github.com/${REPO}/releases/download/${version}/${filename}`;
    
    const binDir = path.join(__dirname, '..', 'bin');
    const tmpFile = path.join(binDir, filename);
    
    // Create bin directory
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir, { recursive: true });
    }
    
    console.log(`   Downloading ${filename}...`);
    await download(downloadUrl, tmpFile);
    
    console.log('   Extracting...');
    if (platform === 'windows') {
      // Use unzip command on Windows
      execSync(`tar -xf "${tmpFile}" -C "${binDir}"`, { stdio: 'ignore' });
    } else {
      await extractTarGz(tmpFile, binDir);
    }
    
    // Rename binary
    const extractedName = `${BINARY_NAME}-${platform}-${arch}${platform === 'windows' ? '.exe' : ''}`;
    const finalName = `${BINARY_NAME}${platform === 'windows' ? '.exe' : ''}`;
    
    const extractedPath = path.join(binDir, extractedName);
    const finalPath = path.join(binDir, finalName);
    
    if (fs.existsSync(extractedPath)) {
      fs.renameSync(extractedPath, finalPath);
    }
    
    // Make executable (Unix)
    if (platform !== 'windows') {
      fs.chmodSync(finalPath, 0o755);
    }
    
    // Cleanup
    fs.unlinkSync(tmpFile);
    
    console.log('✅ ContextPilot installed successfully!');
    console.log('');
    console.log('   Get started:');
    console.log('     cd your-project');
    console.log('     contextpilot init');
    
  } catch (err) {
    console.error('❌ Installation failed:', err.message);
    console.error('');
    console.error('   Try installing manually:');
    console.error('   curl -fsSL https://contextpilot.dev/install.sh | sh');
    process.exit(1);
  }
}

install();
