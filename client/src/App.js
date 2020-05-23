import React from 'react';
import './App.css';
import Editor from './components/Editor/Editor';

function App() {
  return (
    <div className="row">
      <div className="col-md-2">
        <div className="container mt-3">
          <h3>Welcome to the Xlang online editor!</h3>
          <p>
            Feel free to run code, see the examples and see the outputs in the
            terminal!
          </p>
        </div>
      </div>
      <div className="mt-3">
        <Editor />
      </div>
    </div>
  );
}

export default App;
