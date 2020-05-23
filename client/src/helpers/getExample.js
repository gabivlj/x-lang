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
    const uri = url[get];
    const data = await fetch(uri);
    const text = await data.text();

    const run = await fetch('http://localhost:8080/api/v1', {
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
    return '// Your code!';
  }
};

export default getExample;
