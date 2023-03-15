import { render } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import AskALibrarianWidget from './askALibrarianWidget';
import ChatWidget from './ChatWidget';

jest.mock('./askALibrarianWidget');

describe('ChatWidget', () => {
    let mockAskALibrarianWidget;

    beforeEach(() => {
        mockAskALibrarianWidget = {
            run: jest.fn(),
            accessibleSelect: jest.fn(),
            visibleClass: 'showing-chat-frame',
        };
        AskALibrarianWidget.mockImplementation(() => mockAskALibrarianWidget);
    });

    afterEach(() => {
        jest.resetAllMocks();
    });

    test('renders chat widget', () => {
        const { getByText } = render(<ChatWidget />);
        const chatTab = getByText('Chat with us');
        expect(chatTab).toBeInTheDocument();
    });

    test('calls AskALibrarianWidget run method on mount', () => {
        render(<ChatWidget />);
        expect(mockAskALibrarianWidget.run).toHaveBeenCalled();
    });

    test('shows chat frame when chat tab is clicked', () => {
        const { getByText, getByTitle } = render(<ChatWidget />);
        const chatTab = getByText('Chat with us');

        userEvent.click(chatTab);
        const closeChatButton = getByTitle('Close chat window');
        expect(closeChatButton).toBeInTheDocument();
    });

    test('hides chat frame when close button is clicked', () => {
        const { getByText, getByTitle, queryByText } = render(<ChatWidget />);
        const chatTab = getByText('Chat with us');

        userEvent.click(chatTab);
        const closeChatButton = getByTitle('Close chat window');
        userEvent.click(closeChatButton);
        expect(queryByText('Close chat window')).not.toBeInTheDocument();
    });
});
