import React from 'react';
import MonacoEditor from 'react-monaco-editor';

export default function Terminal({ output, error, parse_error }) {
  const outputStr =
    output &&
    output.reduce(
      (prev, message) =>
        `
${prev}
${message.line}: ${message.messages.reduce(
          (prev, now) => `${prev}${now.trim()}`.trim(),
          ``
        )}\n`.trim(),
      `Output:`
    );
  return (
    <MonacoEditor
      value={outputStr || `Xlang is made by Gabriel Villalonga`}
      height="100"
      width="1000"
      options={{
        useTabStops: false,
        readOnly: true,
        tabCompletion: false,
        autoIndent: 'none'
      }}
    />
  );
}
