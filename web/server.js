// require the webpack Node.js library
var webpack = require('webpack');

// Create an instance of the compiler
var compiler = webpack({
  // The first argument is your webpack config
  entry: './entry.js',
  output: {
    path: '.',
    filename: 'bundle.js'
  }
});

function info(err, stats) {
    if (err) {
        console.error(err)
    } else {
        //console.log(stats)
        console.log((stats.endTime-stats.startTime) + " ms")
    }
}

// Run the compiler manually
//compiler.run(info);

// Start watching files and upon change call the callback
compiler.watch(/* watchDelay */ 200, info);
