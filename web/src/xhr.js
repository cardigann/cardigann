
export default function(path) {
  if (window.location.hostname === "localhost" && window.location.port === "3000") {
    return "http://localhost:5060" + path
  }
  return path;
};