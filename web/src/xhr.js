
import parse from 'url-parse';

export default function(path) {
  let url = parse(path, true);

  // rewrite xhr requests to the backend
  if (url.host === "localhost:3000") {
    url.set('port', '5060');
  }

  return url.toString();
};