import { cleanup, render, screen, waitFor } from '@testing-library/react';

import userEvent from '@testing-library/user-event';
import links from '../../testutils/';
import PageFooter from './Footer';

describe('PageFooter', () => {
    let footer;

    beforeEach(() => {
        render(<PageFooter />);
        footer = screen.getByTestId('page-footer');
    });

    afterEach(cleanup);

    test('renders PageFooter component', () => {
        expect(footer).toBeInTheDocument();
    });

    links.forEach((link) => {
        test(`clicking "${link.name}" link should navigate to correct URL`, () => {
            const linkElement = screen.getByText(link.name);
            userEvent.click(linkElement);
            expect(linkElement).toHaveAttribute('href', link.href);

            if (link.target) {
                expect(linkElement).toHaveAttribute('target', link.target);
            }

            if (link.rel) {
                expect(linkElement).toHaveAttribute('rel', link.rel);
            }
        });
    });

    test('renders Twitter logo', async () => {
        const logo = await waitFor(() => screen.getByAltText('Twitter logo').closest('a'));
        expect(logo).toBeInTheDocument();
        expect(logo).toHaveAttribute('href', 'https://twitter.com/nyulibraries');
    });

    test('renders Facebook logo', async () => {
        const logo = await waitFor(() => screen.getByAltText('Facebook logo').closest('a'));
        expect(logo).toBeInTheDocument();
        expect(logo).toHaveAttribute('href', 'https://www.facebook.com/nyulibraries');
    });

    test('renders Instagram logo', async () => {
        const logo = await waitFor(() => screen.getByAltText('Instagram logo').closest('a'));
        expect(logo).toBeInTheDocument();
        expect(logo).toHaveAttribute('href', 'https://www.instagram.com/nyulibraries');
    });

    test('clicking email list subscription link should navigate to correct URL', () => {
        const link = screen.getByText('Subscribe to our email list');
        expect(link).toHaveAttribute('href', 'https://signup.e2ma.net/signup/1934378/1922970/');
    });

    test('renders NYU logo', async () => {
        const logo = await waitFor(() => screen.getByAltText('New York University logo').closest('a'));
        expect(logo).toBeInTheDocument();
        expect(logo).toHaveAttribute('href', 'https://www.nyu.edu');
    });
});

