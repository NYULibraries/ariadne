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
