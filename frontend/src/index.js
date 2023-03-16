import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';


window.IS_PROD_ENV = process.env.NODE_ENV === 'production';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
