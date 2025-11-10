#!/usr/bin/env node
const args = process.argv.slice(2);
const scriptName = args[0] || 'stdio';

async function run() {
  try {
    switch (scriptName) {
      case 'stdio':
        await import('./transport/stdio.js');
        break;
      // case 'sse':
      //   await import('./transport/sse.js');
      //   break;
      // case 'streamableHttp':
      //   await import('./transport/streamableHttp.js');
      //   break;
      default:
        console.error(`Unknown script: ${scriptName}`);
        console.log('Available scripts:');
        console.log('- stdio');
        console.log('- sse');
        console.log('- streamableHttp');
        process.exit(1);
    }
  } catch (error) {
    console.error('Error running script:', error);
    process.exit(1);
  }
}

run();
