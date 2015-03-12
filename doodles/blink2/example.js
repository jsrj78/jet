var start = new Date().getTime();

var OnOff = React.createClass({
  render: function() {
    return (
        <input type="checkbox" />
    )
  }
});
  
var Elapsed = React.createClass({
  render: function() {
    var elapsed = Math.round(this.props.elapsed  / 100);
    var seconds = elapsed / 10 + (elapsed % 10 ? '' : '.0' );
    var message =
      'React has been running for ' + seconds + ' seconds.';

    return <p>{message}</p>;
  }
});

var App = React.createClass({
  render: function() {
    return (
      <OnOff />
      <Elapsed elapsed={new Date().getTime() - start} />
    )
  }
});

setInterval(function() {
  React.render(<App/>, document.getElementById('app'));
}, 50);

// vim: sw=2 ts=2 sts=2
