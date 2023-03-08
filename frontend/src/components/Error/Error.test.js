import { render } from '@testing-library/react';
import Error from './Error';

describe('Error component', () => {
    it('renders error message', () => {
        const { getByText } = render(<Error message="Oops! Something went wrong." />);
        expect(getByText('Oops! Something went wrong.')).toBeInTheDocument();
    });

    it('renders error message', () => {
        const { asFragment } = render(<Error message="Oops! Something went wrong." />);
        expect(asFragment()).toMatchSnapshot();
    });

});
