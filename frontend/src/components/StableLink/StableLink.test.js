import { render, screen } from '@testing-library/react';

import userEvent from '@testing-library/user-event';
import React from 'react';
import StableLink from './StableLink';

describe('StableLink', () => {
    beforeAll(() => {
        Object.defineProperty(navigator, 'clipboard', {
            value: {
                writeText: jest.fn().mockResolvedValue(),
            },
        });
    });

    test('renders the main button', () => {
        render(<StableLink />);
        const mainButton = screen.getByText('🔗 Create a stable link to this page');
        expect(mainButton).toBeInTheDocument();
    });

    test('clicking the main button shows the input field and copy button', () => {
        render(<StableLink />);
        const mainButton = screen.getByText('🔗 Create a stable link to this page');
        userEvent.click(mainButton);
        const inputField = screen.getByRole('textbox');
        const copyButton = screen.getByText('Copy');
        expect(inputField).toBeInTheDocument();
        expect(copyButton).toBeInTheDocument();
    });

    test('clicking the close button hides the input field and copy button', () => {
        render(<StableLink />);
        const mainButton = screen.getByText('🔗 Create a stable link to this page');
        userEvent.click(mainButton);
        const closeButton = screen.getByText('X');
        userEvent.click(closeButton);
        const inputField = screen.queryByRole('textbox');
        const copyButton = screen.queryByText('Copy');
        expect(inputField).not.toBeInTheDocument();
        expect(copyButton).not.toBeInTheDocument();
    });

    test('clicking the copy button copies the link to the clipboard', async () => {
        render(<StableLink />);
        const mainButton = screen.getByText('🔗 Create a stable link to this page');
        userEvent.click(mainButton);
        const copyButton = screen.getByText('Copy');
        userEvent.click(copyButton);
        expect(navigator.clipboard.writeText).toHaveBeenCalledWith(window.location.href);
    });
});
