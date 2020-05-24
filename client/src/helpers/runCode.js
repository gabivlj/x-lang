let uriBuild = `${process.env.REACT_APP_DEPLOY ? 'https' : 'http'}://${process
  .env.REACT_APP_URI || '0.0.0.0'}:${process.env.REACT_APP_PORT ||
  '8080'}/api/v1`;
if (uriBuild.includes('https') && uriBuild.includes(':80')) {
  uriBuild = uriBuild.replace(process.env.REACT_APP_PORT, '');
}
export const uri = uriBuild;

export const runCode = async code => {
  try {
    const run = await fetch(`${uri}`, {
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
