const uri = `http://localhost:8080/api/v1`;

const runCode = async code => {
  try {
    const run = await fetch('http://localhost:8080/api/v1', {
      body: JSON.stringify({ code }),
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      mode: 'cors'
    });
    const json = await run.json();
    return {
      output: json.data.output,
      error: json.data.error,
      errorParsing: json.data.error_parse
    };
  } catch (err) {
    return { error: { line: 0, messages: ['Error running the code'] } };
  }
};

export default runCode;
