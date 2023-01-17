import { render } from '@testing-library/react';
import Citation from './Citation';

describe('Citation', () => {
  it('renders without crashing', () => {
    const metadataPlaceholders = { author: 'John Doe', date: '2023-01-17' };
    render(<Citation metadataPlaceholders={metadataPlaceholders} />);
  });

  it('should pass metadataPlaceholders prop', () => {
    const metadataPlaceholders = { author: 'John Doe', date: '2022-01-17' };
    const { container } = render(<Citation metadataPlaceholders={metadataPlaceholders} />);
    expect(container.querySelector('p').textContent).toBeDefined();
  });

  it('should fail if metadataPlaceholders prop type is not an object', () => {
    const metadataPlaceholders = 'some string';
    expect.assertions(1); // Expect one assertion to be called
    console.error = jest.fn(); // Using jest.fn() to mock the console.error function so that we can check if it is called
    render(<Citation metadataPlaceholders={metadataPlaceholders} />);
    expect(console.error).toHaveBeenCalled(); // Check if the console.error function is called
  });

  it('renders the genre when present', () => {
    const metadataPlaceholders = { genre: 'Research Article' };
    const { getByText } = render(<Citation metadataPlaceholders={metadataPlaceholders} />);
    expect(getByText('Research Article')).toBeInTheDocument();
  });

  it('renders the title when present', () => {
    const metadataPlaceholders = { article_title: 'Example Article Title' };
    const { getByText } = render(<Citation metadataPlaceholders={metadataPlaceholders} />);
    expect(getByText('Example Article Title')).toBeInTheDocument();
  });

  it('renders the author and date when both are present', () => {
    const metadataPlaceholders = { author: 'John Doe', date: '2022-01-01' };
    const { queryByText } = render(<Citation metadataPlaceholders={metadataPlaceholders} />);
    expect(queryByText(/John Doe.*2022-01-01/)).toBeInTheDocument();
  });

  it('renders the ISSN when present', () => {
    const metadataPlaceholders = { issn: '1234-5678' };
    const { getByText } = render(<Citation metadataPlaceholders={metadataPlaceholders} />);
    expect(getByText('ISSN:')).toBeInTheDocument();
    expect(getByText('1234-5678')).toBeInTheDocument();
  });

  it('should not display any empty values in the rendered output', () => {
    const metadataPlaceholders = {
      genre: 'Article',
      date: '1992',
      article_title: 'Test Article Title',
      journal_title: '',
      volume: '',
      issue: '',
      start_page: '',
      end_page: '206',
      issn: '',
      author: 'John Don',
    };
    const { queryByText } = render(<Citation metadataPlaceholders={metadataPlaceholders} />);
    expect(queryByText('Article')).toBeInTheDocument();
    expect(queryByText(/1992/)).toBeInTheDocument();
    expect(queryByText('Published in Journal')).not.toBeInTheDocument();
    expect(queryByText('Test Article Title')).toBeInTheDocument();
    expect(queryByText('Volume ')).toBeNull();
    expect(queryByText('Issue ')).toBeNull();
    expect(queryByText('Page ')).toBeNull();
    expect(queryByText('ISSN:')).toBeNull();
  });
});
