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
    this.interval = setInterval(function() {
      if (this.state.enabled) {
        this.setState({blink: !this.state.blink});
      }
    }.bind(this), 500);
  },

  componentWillUnmount: function() {
    clearInterval(this.interval);
  },

  handleClick: function(event) {
    this.setState({enabled: !this.state.enabled});
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
