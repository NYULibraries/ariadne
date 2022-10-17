import { render, screen, waitFor, waitForElementToBeRemoved } from '@testing-library/react';
import '@testing-library/jest-dom';
import List from './List';

const queryString =
  '?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat';

const searchURL = `${process.env.REACT_APP_API_URL}` + queryString;

beforeEach(() => {
  delete window.location;
  window.location = new URL(searchURL);
});

afterEach(() => {
  delete window.location;
  window.location = new URL('http://localhost:3000');
});

test('renders with a className of list-group', () => {
  const { container } = render(<List />);
  expect(container.getElementsByClassName('list-group').length).toBe(1);
});

test('renders correctly', () => {
  const { container } = render(<List />);
  expect(container.firstChild).toMatchSnapshot();
});

test('renders an error when error', () => {
  const { container } = render(<List error="Error" />);
  expect(container.firstChild).toMatchSnapshot();
});

test('renders a list of elements', () => {
  const { container } = render(<List elements={[{ id: 1 }, { id: 2 }]} />);
  expect(container.firstChild).toMatchSnapshot();
});

test('renders a last element', () => {
  const { container } = render(<List lastElement={{ id: 1 }} />);
  expect(container.firstChild).toMatchSnapshot();
});

test('renders a list of elements and a last element', () => {
  const { container } = render(<List elements={[{ id: 1 }, { id: 2 }]} lastElement={{ id: 3 }} />);
  expect(container.firstChild).toMatchSnapshot();
});

test('renders Loading...', () => {
  render(<List />);
  const linkElement = screen.getByText(/Loading/i);
  expect(linkElement).toBeInTheDocument();
});

test('renders a E Journal Full Text link', async () => {
  render(<List />);
  const linkElement = await waitFor(() => screen.getByText(/E Journal Full Text/i));
  expect(linkElement).toBeInTheDocument();
});

test('renders a Gale General OneFile link', async () => {
  render(<List />);
  const linkElement = await waitFor(() => screen.getByText(/Gale General OneFile/i));
  expect(linkElement).toBeInTheDocument();
});

test('renders Ask a Librarian', async () => {
  render(<List />);
  const linkElement = await waitFor(() => screen.getByText(/Ask a Librarian/i));
  expect(linkElement).toBeInTheDocument();
});

test('Loading... no longer present in the DOM after loading data', async () => {
  const { getByText } = render(<List />);
  await waitForElementToBeRemoved(() => getByText(/Loading/i));
});
