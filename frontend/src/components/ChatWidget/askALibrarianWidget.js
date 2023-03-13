const CLICK_EVENT = 'click';
const KEYDOWN_EVENT = 'keydown';
const SPACE_KEY = 'Space';
const ENTER_KEY = 'Enter';

function accessibleSelect(event) {
    if (event.type === CLICK_EVENT) {
        return true;
    } else if (event.type === KEYDOWN_EVENT) {
        const code = event.key || event.code;
        if (code === SPACE_KEY || code === ENTER_KEY) {
            return true;
        }
    }
    return false;
}

export default class AskALibrarianWidget {
    constructor() {
        this.selector = document.querySelectorAll('.js-toggle-chat, .js-toggle-chat-from-link');
        this.triggerOnEvents = [CLICK_EVENT, KEYDOWN_EVENT];
        this.visibleClass = 'showing-chat-frame';
        this.linkSelectorClass = 'js-toggle-chat-from-link';
        this.chatWidgetFocusSelector = ['aside', '#chat_widget'];
    }

    run() {
        const obj = this;
        obj.selector.forEach(link => {
            obj.triggerOnEvents.forEach(event => {
                link.addEventListener(event, function (event) {
                    if (accessibleSelect(event)) {
                        event.preventDefault();
                        document.body.classList.toggle(obj.visibleClass);
                        // If triggered from another link on the page focus on the widget
                        if (this.classList.contains(obj.linkSelectorClass) && document.body.classList.contains(obj.visibleClass)) {
                            document.querySelector(obj.chatWidgetFocusSelector).focus();
                        }
                    }
                });
            });
        });
    }
}


