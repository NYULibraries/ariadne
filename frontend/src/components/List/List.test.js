import '@testing-library/jest-dom';

import { render, screen } from '@testing-library/react';

import List from './List';

const mockLinks = [
  {
    display_name: 'Example Website',
    url: 'https://www.example.com',
  },
  {
    display_name: 'Example Organization',
    url: 'https://www.example.org',
  }
];

describe('List', () => {
  it('renders without error when links prop is passed', () => {
    render(<List links={mockLinks} />);
    const linkElements = screen.getAllByRole('link');
    expect(linkElements.length).toBe(2);
    expect(linkElements[0]).toHaveAttribute('href', 'https://www.example.com');
    expect(linkElements[1]).toHaveAttribute('href', 'https://www.example.org');
  });

  it('renders "No results found" message when links prop is empty', () => {
    render(<List links={[]} />);
    const noResultsMessage = screen.getByText('No results found');
    expect(noResultsMessage).toBeInTheDocument();
  });

  it('renders error message when error prop is passed', () => {
    const errorMessage = 'Failed to load links';
    render(<List error={errorMessage} />);
    const errorElement = screen.getByText(errorMessage);
    expect(errorElement).toBeInTheDocument();
  });

  it('renders loading indicator when loading prop is true', () => {
    render(<List loading={true} />);
    const loadingElement = screen.getByLabelText('Loading...');
    expect(loadingElement).toBeInTheDocument();
  });

  it('matches snapshot when links prop is passed', () => {
    const { container } = render(<List links={mockLinks} />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it('matches snapshot when links prop is empty', () => {
    const { container } = render(<List links={[]} />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it('matches snapshot when error prop is passed', () => {
    const errorMessage = 'Failed to load links';
    const { container } = render(<List error={errorMessage} />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it('matches snapshot when loading prop is true', () => {
    const { container } = render(<List loading={true} />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it('renders the correct link text', () => {
    render(<List links={mockLinks} />);
    const linkElements = screen.getAllByRole('link');
    expect(linkElements.length).toBe(2);
    expect(linkElements[0]).toHaveTextContent('Example Website');
    expect(linkElements[1]).toHaveTextContent('Example Organization');
  });

  it('renders an Error component with the correct message when error prop is passed', () => {
    const errorMessage = 'Failed to load links';
    render(<List error={errorMessage} />);
    const errorComponent = screen.getByRole('alert');
    expect(errorComponent).toBeInTheDocument();
    expect(errorComponent).toHaveTextContent(errorMessage);
  });

  it('renders a Loader component when loading prop is true', () => {
    render(<List loading={true} />);
    const loaderComponent = screen.getByLabelText('Loading...');
    expect(loaderComponent).toBeInTheDocument();
  });
});
