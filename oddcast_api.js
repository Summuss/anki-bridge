const {buildUrl, voices} = require('E:\\lib\\js\\oddcast-tts-demo');

if (process.argv.length < 4) {
    console.log('args length less than 2');
    process.exit(1);
}
text = process.argv[2]
voice = process.argv[3]

if (!text) {
    console.log('text is empty');
    process.exit(1)
}
if (!voice) {
    console.log('voice is empty')
    process.exit(1)
}

url = buildUrl(text, eval(`voices.${voice}`));
console.log(url)
process.exit(0)