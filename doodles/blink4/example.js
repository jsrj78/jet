var OnOff = React.createClass({
  render: function() {
    var onClick = this.props.onClick;
    var value = this.props.value;
    return (
      <div>
        <input type="checkbox" id="onoff" checked={value} onClick={onClick} />
        <label htmlFor="onoff">Enable</label>
      </div>
    );
  }
});

var Indicator = React.createClass({
  render: function() {
    var style = {width: 20, height: 20, border: "1px solid black"};
    style.background = this.props.on ? "red" : "white";
    return (
      <div style={style} />
    );
  }
});

var App = React.createClass({
  getInitialState: function() {
    return { enabled: false, blink: false };
  },

  componentWillMount: function() {
    this.ws = new WebSocket("ws://localhost:8000/ws");
    this.ws.onmessage = function(event) {
      this.setState(JSON.parse(event.data));
    }.bind(this);
  },

  componentWillUnmount: function() {
    this.ws.close();
  },

  handleClick: function(event) {
    var enabled = !this.state.enabled;
    this.ws.send(JSON.stringify({enabled: enabled}));
    this.setState({enabled: enabled});
  },

  render: function() {
    return (
      <div>
        <p><OnOff value={this.state.enabled} onClick={this.handleClick} /></p>
        <p><Indicator on={this.state.blink} /></p>
      </div>
    );
  }
});

React.render(<App/>, document.getElementById('app'));
