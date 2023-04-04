import AskLibrarian, { ASK_LIBRARIAN_TEXT, ASK_LIBRARIAN_URL } from './AskLibrarian';

import { render } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

describe('AskLibrarian component', () => {
    it('renders "Need Help?', () => {
        const { getByText } = render(<AskLibrarian />);
        expect(getByText('Need Help?')).toBeInTheDocument();
    });

    it('renders the correct ASK_LIBRARIAN_TEXT and ASK_LIBRARIAN_URL for the second link', () => {
        const { getAllByText } = render(<AskLibrarian />);
        const links = getAllByText(ASK_LIBRARIAN_TEXT);
        const secondLink = links[1];
        expect(secondLink).toBeInTheDocument();
        expect(secondLink.getAttribute('href')).toBe(ASK_LIBRARIAN_URL);
    });


    it('renders additional resources', () => {
        const { queryAllByText } = render(<AskLibrarian />);
        expect(queryAllByText(/Additional Resources/).length).toBeGreaterThan(0);
    });


    it('opens the ASK_A_LIBRARIAN_URL in a new tab when the second link is clicked', () => {
        const { getAllByText } = render(<AskLibrarian />);
        const links = getAllByText(ASK_LIBRARIAN_TEXT);
        const secondLink = links[1];
        userEvent.click(secondLink);
        expect(secondLink).toHaveAttribute('target', '_blank');
        expect(secondLink).toHaveAttribute('rel', 'noreferrer');
        expect(secondLink).toHaveAttribute('href', ASK_LIBRARIAN_URL);
    });

    it('renders both Ask a Librarian elements as visible', () => {
        const { getAllByText } = render(<AskLibrarian />);
        const links = getAllByText(ASK_LIBRARIAN_TEXT);
        expect(links.length).toBe(2);

        links.forEach((link) => {
            expect(link).toBeVisible();
        });
    });

    it('matches snapshot', () => {
        const { container } = render(<AskLibrarian />);
        expect(container.firstChild).toMatchSnapshot();
    });
});
