import React, { useEffect } from "react";

import AskALibrarianWidget from "./askALibrarianWidget";

const ChatWidget = () => {
    useEffect(() => {
        const chatWidget = new AskALibrarianWidget();
        chatWidget.run();
    }, []);

    return (
        <aside id="chat_widget" tabIndex="-1">
            <button className="button chat-tab ss-chat js-toggle-chat">
                Chat with us
            </button>
            <div className="chat-frame-wrap">
                <button
                    className="chat-close ss-icon js-toggle-chat"
                    title="Close chat window"
                    aria-label="Close chat window"
                >
                    X
                </button>
                <div
                    className="libraryh3lp"
                    data-lh3-jid="nyu-aal-chat@chat.libraryh3lp.com"
                >
                    <iframe
                        title="Ask a Librarian chat"
                        src="https://libraryh3lp.com/chat/nyu-aal-chat@chat.libraryh3lp.com?skin=23106&amp;referer=https%3A%2F%2Flibrary.nyu.edu%2F"
                        style={{
                            width: "186px",
                            height: "250px",
                            border: "2px solid rgb(192, 192, 192)",
                        }}
                    ></iframe>
                </div>
            </div>
        </aside>
    );
};

export default ChatWidget;
