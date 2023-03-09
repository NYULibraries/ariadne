import '@testing-library/jest-dom';

import { render, screen } from '@testing-library/react';

import List from './List';

const mockLinks = [
  {
    target_url: 'https://www.example.com',
    target_public_name: 'Example Website'
  },
  {
    target_url: 'https://www.example.org',
    target_public_name: 'Example Organization'
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

});
