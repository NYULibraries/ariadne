import { act, render, screen } from '@testing-library/react';

import React from 'react';
import StableLink from './StableLink';
import userEvent from '@testing-library/user-event';

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
        const mainButton = screen.getByText('Copy a stable link to this page');
        expect(mainButton).toBeInTheDocument();
    });

    test('clicking the main button shows "Copied!" message', async () => {
        render(<StableLink />);
        const mainButton = screen.getByText('Copy a stable link to this page');

        await act(async () => {
            userEvent.click(mainButton);
        });

        expect(screen.getByText('Copied!')).toBeInTheDocument();
    });

    test('clicking the main button calls navigator.clipboard.writeText', async () => {
        render(<StableLink />);
        const mainButton = screen.getByText('Copy a stable link to this page');

        await act(async () => {
            userEvent.click(mainButton);
        });
        expect(navigator.clipboard.writeText).toHaveBeenCalledWith(window.location.href);
    });
});
