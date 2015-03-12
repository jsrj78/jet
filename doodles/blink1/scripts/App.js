import React from 'react';

class OnOff extends React.Component {
  render() {
    return (
      <input type="checkbox" id="onoff" />
      <label for="onoff">Enable</label>
    );
  }
}

export default class App extends React.Component {
  render() {
    return (
      <h1>Hello, world.</h1>
      <OnOff />
    );
  }
}

// vim: sw=2 ts=2 sts=2
