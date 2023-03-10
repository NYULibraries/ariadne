import AskLibrarian, { ASK_LIBRARIAN_TEXT, ASK_LIBRARIAN_URL } from './AskLibrarian';

import { render } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

describe('AskLibrarian component', () => {
    it('renders "Need Help?" when not loading', () => {
        const { getByText } = render(<AskLibrarian />);
        expect(getByText('Need Help?')).toBeInTheDocument();
    });

    it('renders the correct ASK_LIBRARIAN_TEXT and ASK_LIBRARIAN_URL', () => {
        const { getByText } = render(<AskLibrarian />);
        const link = getByText(ASK_LIBRARIAN_TEXT);
        expect(link).toBeInTheDocument();
        expect(link.getAttribute('href')).toBe(ASK_LIBRARIAN_URL);
    });

    it('renders additional resources when not loading', () => {
        const { queryAllByText } = render(<AskLibrarian />);
        expect(queryAllByText(/Additional Resources/).length).toBeGreaterThan(0);
    });


    it('opens the ASK_A_LIBRARIAN_URL in a new tab when the link is clicked', () => {
        const { getByText } = render(<AskLibrarian />);
        const link = getByText(ASK_LIBRARIAN_TEXT);
        userEvent.click(link);
        expect(link).toHaveAttribute('target', '_blank');
        expect(link).toHaveAttribute('rel', 'noopener noreferrer');
        expect(link).toHaveAttribute('href', ASK_LIBRARIAN_URL);
    });

    it('matches snapshot', () => {
        const { container } = render(<AskLibrarian />);
        expect(container.firstChild).toMatchSnapshot();
    });
});
