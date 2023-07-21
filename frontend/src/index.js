import './index.css';

import App from './App';
import React from 'react';
import ReactDOM from 'react-dom/client';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

// statuspage embed
(function () {
  var env = process.env.REACT_APP_ENV;
  if (env === 'production') {
    var script = document.createElement('script');
    script.src = 'https://cdn.library.nyu.edu/statuspage-embed/index.min.js';
    script.async = true;
    document.body.appendChild(script);
  }
})();

// libraryh3lp chat widget
(function () {
  var x = document.createElement("script"); x.type = "text/javascript"; x.async = true;
  x.src = (document.location.protocol === "https:" ? "https://" : "http://") + "libraryh3lp.com/js/libraryh3lp.js?7516";
  var y = document.getElementsByTagName("script")[0]; y.parentNode.insertBefore(x, y);
})();

// chatwidget-embed
(function () {
  var s = document.createElement('script');
  s.type = 'text/javascript';
  s.async = true;
  s.src = 'https://cdn.library.nyu.edu/chatwidget-embed/index.min.js';
  var x = document.getElementsByTagName('script')[0];
  x.parentNode.insertBefore(s, x);
})();