import './ChatWidget.css';

import { renderChatWidget } from './ChatWidget.js';

// Create a container for the chat widget
const app = document.querySelector('#app');
const chatContainer = document.createElement('div');
chatContainer.className = 'chat-container';
app.appendChild(chatContainer);

// Call the chat widget function to render the chat widget
renderChatWidget();
