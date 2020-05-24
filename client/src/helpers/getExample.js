import { uri } from './runCode';

const host = 'https://raw.githubusercontent.com/gabivlj/x-lang/master/examples';

const url = {
  Arrays: `${host}/arrays.xlang`,
  Integers: `${host}/integers.xlang`,
  Strings: `${host}/strings.xlang`,
  Hashmaps: `${host}/hashtable_example.xlang`
};
const getExample = async get => {
  try {
    if (!url[get]) return '// Your code!';
    const URIGithub = url[get];
    const data = await fetch(URIGithub);
    const text = await data.text();

    const run = await fetch(uri, {
      body: JSON.stringify({ code: text }),
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      mode: 'cors'
    });
    const json = await run.json();
    return [text, json];
  } catch (err) {
    return [
      '// Your code',
      { data: { output: { message: { messages: [''], line: '0' } } } }
    ];
  }
};

export default getExample;
