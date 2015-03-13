import React from 'react';

class OnOff extends React.Component {
  render() {
    return (
      <div>
        <input type="checkbox" id="onoff" />
        <label htmlFor="onoff">Enable</label>
      </div>
    );
  }
}

export default class App extends React.Component {
  render() {
    return (
      <div>
        <h1>JET Blink 1</h1>
        <OnOff />
      </div>
    );
  }
}
