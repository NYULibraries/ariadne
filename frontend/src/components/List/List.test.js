import { render } from '@testing-library/react';
import List from './List';

test('renders with a className of list-group', () => {
  const { container } = render(<List />);
  expect(container.getElementsByClassName('list-group').length).toBe(1);
});

test('renders correctly', () => {
  const { container } = render(<List />);
  expect(container.firstChild).toMatchSnapshot();
});

test('renders Loading... when loading', () => {
  const { container } = render(<List loading />);
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
