import React, { useState } from 'react';
import MonacoEditor from 'react-monaco-editor';
import getExample from '../../helpers/getExample';
import runCode from '../../helpers/runCode';
import Terminal from './Terminal';

export default function Editor() {
  const [code, setCode] = useState('// Write some xlang code...');
  const [terminalObject, setTerminalObject] = useState({});

  return (
    <div style={{ borderRadius: '1px', borderLeft: 'solid' }}>
      <select
        className="m-3"
        onChange={e => {
          setCode('// Loading example from Github...');
          getExample(e.target.value).then(([text, terminal]) => {
            setCode(text);

            setTerminalObject({
              output: terminal.data.output,
              error: terminal.data.error
            });
          });
        }}
      >
        <option>Your code!</option>
        <option>Integers</option>
        <option>Strings</option>
        <option>Arrays</option>
        <option>Hashmaps</option>
      </select>
      <button
        className="btn btn-primary"
        onClick={e => {
          e.preventDefault();
          runCode(code).then(res => setTerminalObject(res));
        }}
      >
        Run
      </button>
      <div className="pl-3">
        <MonacoEditor
          width="1300"
          height="700"
          language="rust"
          theme="vs-dark"
          onChange={e => {
            setCode(e);
          }}
          options={{ fontSize: 16 }}
          value={code}
        ></MonacoEditor>
      </div>

      <div className="pl-3 mt-3">
        <div>
          <Terminal output={terminalObject.output}></Terminal>
        </div>
      </div>
    </div>
  );
}
