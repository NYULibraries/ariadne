import './ChatWidget.css';

import AskALibrarianWidget from './askALibrarianWidget';

let showChat = false;

const handleChatToggle = () => {
    showChat = !showChat;
    renderChatWidget();
};

const chatWidget = new AskALibrarianWidget();
chatWidget.run();

export const renderChatWidget = () => {
    const chatContainer = document.querySelector('.chat-container');
    chatContainer.id = 'chat_widget';
    chatContainer.innerHTML = '';

    if (!showChat) {
        const button = document.createElement('button');
        button.className = 'button ss-chat chat-button';
        button.onclick = handleChatToggle;

        const img = document.createElement('img');
        img.src = '/chat-icon.svg';
        img.alt = 'Chat-Symbol.';
        img.className = 'chat-icon';

        button.appendChild(img);
        button.appendChild(document.createTextNode('Chat with us'));

        chatContainer.appendChild(button);
    } else {
        const closeButton = document.createElement('button');
        closeButton.className = 'chat-close ss-icon js-toggle-chat';
        closeButton.title = 'Close chat window';
        closeButton.setAttribute('aria-label', 'Close chat window');
        closeButton.onclick = handleChatToggle;
        closeButton.textContent = '\u00D7';

        const libraryh3lpDiv = document.createElement('div');
        libraryh3lpDiv.className = 'libraryh3lp';
        libraryh3lpDiv.setAttribute('data-lh3-jid', 'nyu-aal-chat@chat.libraryh3lp.com');

        const iframe = document.createElement('iframe');
        iframe.title = 'Ask a Librarian chat';
        iframe.src = 'https://libraryh3lp.com/chat/nyu-aal-chat@chat.libraryh3lp.com?skin=23106&amp;referer=https%3A%2F%2Flibrary.nyu.edu%2F';

        libraryh3lpDiv.appendChild(iframe);

        const chatFrameWrap = document.createElement('div');
        chatFrameWrap.className = 'chat-frame-wrap';

        chatFrameWrap.appendChild(closeButton);
        chatFrameWrap.appendChild(libraryh3lpDiv);

        chatContainer.appendChild(chatFrameWrap);
    }
};

