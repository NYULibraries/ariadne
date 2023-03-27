import "./ChatWidget.css";

import React, { useState } from "react";

import AskALibrarianWidget from "./askALibrarianWidget";

const ChatWidget = () => {
    const chatWidget = new AskALibrarianWidget();
    chatWidget.run();

    const [showChat, setShowChat] = useState(false);

    const handleChatToggle = () => {
        setShowChat(!showChat);
    };

    return (
        <>
            <aside id="chat_widget" tabIndex="-1">
                <div className="chat-container">
                    {!showChat && (
                        <button className="button chat-tab ss-chat chat-button" onClick={handleChatToggle}>
                            Chat with us
                        </button>
                    )}
                    <div className="chat-frame-wrap">
                        {showChat && (
                            <>
                                <button
                                    className="chat-close ss-icon js-toggle-chat"
                                    title="Close chat window"
                                    aria-label="Close chat window"
                                    onClick={handleChatToggle}
                                >
                                    &times;
                                </button>
                                <div className="libraryh3lp" data-lh3-jid="nyu-aal-chat@chat.libraryh3lp.com">
                                    <iframe
                                        title="Ask a Librarian chat"
                                        src="https://libraryh3lp.com/chat/nyu-aal-chat@chat.libraryh3lp.com?skin=23106&amp;referer=https%3A%2F%2Flibrary.nyu.edu%2F"
                                    ></iframe>
                                </div>
                            </>
                        )}
                    </div>
                </div>
            </aside>
        </>
    );
};

export default ChatWidget;
