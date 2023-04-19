import { act, render, screen } from '@testing-library/react';

import userEvent from '@testing-library/user-event';
import React from 'react';
import StableLink from './StableLink';

describe('StableLink', () => {
    beforeAll(() => {
        Object.defineProperty(navigator, 'clipboard', {
            value: {
                writeText: jest.fn(),
            },
        });
        // To suppress this console error output during testing,mock the console.error function, preventing it from outputting the error message during the test
        jest.spyOn(console, 'error').mockImplementation(() => { });
    });

    afterAll(() => {
        // Restore the original implementation of console.error after all the tests have run
        console.error.mockRestore();
    });

    beforeEach(() => {
        jest.resetAllMocks();
    });

    test('renders the main button', () => {
        render(<StableLink />);
        const mainButton = screen.getByText('Copy a stable link to this page');
        expect(mainButton).toBeInTheDocument();
    });

    test('clicking the main button shows "Stable link copied to clipboard!" message', async () => {
        render(<StableLink />);
        const mainButton = screen.getByText('Copy a stable link to this page');

        await act(async () => {
            userEvent.click(mainButton);
        });

        expect(screen.getByText('Stable link copied to clipboard!')).toBeInTheDocument();
    });

    test('clicking the main button calls navigator.clipboard.writeText', async () => {
        render(<StableLink />);
        const mainButton = screen.getByText('Copy a stable link to this page');

        await act(async () => {
            userEvent.click(mainButton);
        });
        expect(navigator.clipboard.writeText).toHaveBeenCalledWith(window.location.href);
    });

    test('displays error message when copying fails', async () => {
        navigator.clipboard.writeText.mockRejectedValue(new Error('Error copying text'));

        render(<StableLink />);
        const mainButton = screen.getByText('Copy a stable link to this page');

        await act(async () => {
            userEvent.click(mainButton);
        });

        const errorMessages = screen.getAllByText('Error copying link. Please try again or manually copy the link.');
        expect(errorMessages.length).toBeGreaterThanOrEqual(1);
    });
});
