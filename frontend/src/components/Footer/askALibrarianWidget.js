import $ from 'jquery';

export default class askALibrarianWidget {
    constructor() {
        this.$selector = $('.js-toggle-chat, .js-toggle-chat-from-link');
        this.triggerOnEvents = 'click keydown';
        this.visibleClass = 'showing-chat-frame';
        this.linkSelectorClass = 'js-toggle-chat-from-link';
        this.chatWidgetFocusSelector = 'aside#chat_widget';
    }

    run = () => {
        const obj = this;
        obj.$selector.on(obj.triggerOnEvents, function (event) {
            if (accessibleSelect(event)) {
                event.preventDefault();
                document.body.classList.toggle(obj.visibleClass);
                // If triggered from another link on the page focus on the widget
                if (this.classList.contains(obj.linkSelectorClass) && document.body.classList.contains(obj.visibleClass)) {
                    $(obj.chatWidgetFocusSelector).trigger();
                }
            }
        });
    }
}

// An accessible select is a mouse click, a space bar (32), or a return (13)
export function accessibleSelect(event) {
    if (event.type === 'click') {
        return true;
    } else if (event.type === 'keydown') {
        const code = event.charCode || event.keyCode;
        if ((code === 32) || (code === 13)) {
            return true;
        }
    } else {
        return false;
    }
}

// Toggle `aria-expanded` accessibly attributes based on value of visible
export function toggleAriaExpanded($link, visible) {
    if (visible) {
        $link.attr('aria-expanded', 'true');
    } else {
        $link.attr('aria-expanded', 'false');
    }
}