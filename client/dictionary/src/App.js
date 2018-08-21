import React, { Component } from 'react';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <h1 className="App-title">React client for distributed dictionary application in Go</h1>
        </header>
        <p className="App-intro">
          Type in an English term to get the definition in English
        </p>
      </div>
    );
  }
}

export default App;
