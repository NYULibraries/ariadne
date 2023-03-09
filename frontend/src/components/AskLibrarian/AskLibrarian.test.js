import AskLibrarian, { ASK_LIBRARIAN_TEXT, ASK_LIBRARIAN_URL } from './AskLibrarian';

import { render } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

describe('AskLibrarian component', () => {
    it('renders "Need Help?" when not loading', () => {
        const { getByText } = render(<AskLibrarian loading={false} />);
        expect(getByText('Need Help?')).toBeInTheDocument();
    });

    it('renders the correct ASK_LIBRARIAN_TEXT and ASK_LIBRARIAN_URL', () => {
        const { getByText } = render(<AskLibrarian loading={false} />);
        const link = getByText(ASK_LIBRARIAN_TEXT);
        expect(link).toBeInTheDocument();
        expect(link.getAttribute('href')).toBe(ASK_LIBRARIAN_URL);
    });

    it('renders additional resources when not loading', () => {
        const { queryAllByText } = render(<AskLibrarian loading={false} />);
        expect(queryAllByText(/Additional Resources/).length).toBeGreaterThan(0);
        // expect(queryAllByText(/Visit our online tutorials for tips on searching the catalog and getting library resources/).length).toBeGreaterThan(0);
        // expect(queryAllByText('Discover subject specific resources using expert curated research guides').length).toBeGreaterThan(0);
        // expect(queryAllByText('Explore the complete list of library services').length).toBeGreaterThan(0);
        // expect(queryAllByText('Reach out to the Libraries on our Instagram').length).toBeGreaterThan(0);
        // expect(queryAllByText('Search WorldCat for items in nearby libraries').length).toBeGreaterThan(0);
    });

    it('does not render anything when loading', () => {
        const { container } = render(<AskLibrarian loading={true} />);
        expect(container.firstChild).toBeNull();
    });

    it('opens the ASK_A_LIBRARIAN_URL in a new tab when the link is clicked', () => {
        const { getByText } = render(<AskLibrarian loading={false} />);
        const link = getByText(ASK_LIBRARIAN_TEXT);
        userEvent.click(link);
        expect(link).toHaveAttribute('target', '_blank');
        expect(link).toHaveAttribute('rel', 'noopener noreferrer');
        expect(link).toHaveAttribute('href', ASK_LIBRARIAN_URL);
    });

    it('matches snapshot when not loading', () => {
        const { container } = render(<AskLibrarian loading={false} />);
        expect(container.firstChild).toMatchSnapshot();
    });

    it('matches snapshot when loading', () => {
        const { container } = render(<AskLibrarian loading={true} />);
        expect(container.firstChild).toMatchSnapshot();
    });
});
