import { render } from '@testing-library/react';
import Citation from './Citation';

describe('Citation', () => {
  test('renders without crashing', () => {
    render(<Citation />);
  });

  test('renders title', () => {
    const { getAllByText } = render(<Citation />);
    const titles = getAllByText((_content, element) =>
      element.textContent.includes('Informed consent and the geriatric dental patient')
    );
    expect(titles.length).toBeGreaterThan(0);
  });

  test('renders author and date', () => {
    const { getAllByText } = render(<Citation />);
    const authors = getAllByText((_content, element) => element.textContent.includes('Odom, John G'));
    const dates = getAllByText((_content, element) => element.textContent.includes('1992'));
    expect(authors.length).toBeGreaterThan(0);
    expect(dates.length).toBeGreaterThan(0);
  });

  test('renders ISSN', () => {
    const { getByText } = render(<Citation />);
    const issn = getByText(/issn/i);
    expect(issn).toBeInTheDocument();
  });

  test('renders journal information', () => {
    const { getAllByText } = render(<Citation />);
    const journalTitles = getAllByText((_content, element) =>
      element.textContent.includes('Special Care in Dentistry')
    );
    const volumes = getAllByText((_content, element) => element.textContent.includes('Volume 12'));
    const issues = getAllByText((_content, element) => element.textContent.includes('Issue 5'));
    const startPages = getAllByText((_content, element) => element.textContent.includes('Page 202'));
    const endPages = getAllByText((_content, element) => element.textContent.includes('-206'));
    expect(journalTitles.length).toBeGreaterThan(0);
    expect(volumes.length).toBeGreaterThan(0);
    expect(issues.length).toBeGreaterThan(0);
    expect(startPages.length).toBeGreaterThan(0);
    expect(endPages.length).toBeGreaterThan(0);
  });
});
