import { render } from '@testing-library/react';
import Citation from './Citation';

describe('Citation', () => {
  it('renders the genre when present', () => {
    window.history.pushState({}, null, '?genre=journal');
    const { queryByText } = render(<Citation />);

    expect(queryByText('Journal')).toBeInTheDocument();
  });

  it('renders the genre when present with rft prefix', () => {
    window.history.pushState({}, null, '?rft.genre=article');
    const { queryByText } = render(<Citation />);

    expect(queryByText('Article')).toBeInTheDocument();
  });

  it('renders the genre when present for book chapters', () => {
    window.history.pushState({}, null, '?genre=bookitem');
    const { queryByText } = render(<Citation />);

    expect(queryByText('Book Chapter')).toBeInTheDocument();
  });

  it('renders the genre when present and uppercased', () => {
    window.history.pushState({}, null, '?genre=PrEPrINt');
    const { queryByText } = render(<Citation />);

    expect(queryByText('Pre-print')).toBeInTheDocument();
  });

  it('renders the title when present for an article', () => {
    window.history.pushState({}, null, '?atitle=Example+Article+Title');
    const { queryByText } = render(<Citation />);

    expect(queryByText('Example Article Title')).toBeInTheDocument();
  });

  it('renders the title when present for a book', () => {
    window.history.pushState({}, null, '?btitle=Sketching+user+experiences+:+getting+the+design+right+and+the+right+design');
    const { queryByText } = render(<Citation />);

    expect(queryByText('Sketching user experiences : getting the design right and the right design')).toBeInTheDocument();
  });

  it('renders the title and container title when present for a book chapter', () => {
    window.history.pushState({}, null, '?atitle=Theorizing Matriarchy in Africa: Kinship Ideologies and Systems in Africa and Europe.&aulast=Amadiume&aufirst=Ifi&title=Re-inventing Africa: Matriarchy, Religion, and Culture');
    const { queryByText } = render(<Citation />);

    expect(queryByText('Theorizing Matriarchy in Africa: Kinship Ideologies and Systems in Africa and Europe.')).toBeInTheDocument();
    expect(queryByText('Published in Re-inventing Africa : matriarchy, religion, and culture')).toBeInTheDocument();
  });

  it('renders the author and date when both are present', () => {
    window.history.pushState({}, null, '?aulast=Doe&aufirst=John&date=1999');
    const { queryByText } = render(<Citation />);

    expect(queryByText(/Doe, John.*1999/)).toBeInTheDocument();
  });

  it('renders the ISSN when present', () => {
    window.history.pushState({}, null, '?issn=1234-5678');
    const { queryByText } = render(<Citation />);

    expect(queryByText('ISSN:')).toBeInTheDocument();
    expect(queryByText('1234-5678')).toBeInTheDocument();
  });

  it('renders the ISBN when present', () => {
    window.history.pushState({}, null, '?rft.isbn=9780080552903');
    const { queryByText } = render(<Citation />);

    expect(queryByText('ISBN:')).toBeInTheDocument();
    expect(queryByText('9780080552903')).toBeInTheDocument();
  });

  it('does not render ISBN when empty', () => {
    window.history.pushState({}, null, '?isbn=&title=something');
    const { queryByText } = render(<Citation />);

    expect(queryByText('ISBN:')).not.toBeInTheDocument();
  });

  it('renders the publisher when present', () => {
    window.history.pushState({}, null, '?rft.pub=Elsevier/Morgan+Kaufmann');
    const { queryByText } = render(<Citation />);

    expect(queryByText('Publisher:')).toBeInTheDocument();
    expect(queryByText('Elsevier/Morgan Kaufmann')).toBeInTheDocument();
  });

  it('should not display any empty values in the rendered output', () => {

    window.history.pushState({}, null, '?genre=Article&atitle=Test+Article+Title&date=1992');
    const { queryByText } = render(<Citation />);

    expect(queryByText('Article')).toBeInTheDocument();
    expect(queryByText(/1992/)).toBeInTheDocument();
    expect(queryByText('Published in')).not.toBeInTheDocument();
    expect(queryByText('Test Article Title')).toBeInTheDocument();
    expect(queryByText('Volume ')).toBeNull();
    expect(queryByText('Issue ')).toBeNull();
    expect(queryByText('Page ')).toBeNull();
    expect(queryByText('ISSN:')).toBeNull();
  });
});
