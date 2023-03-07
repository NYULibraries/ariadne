import Loader, { LOADING_TEXT } from './Loader';

import { render } from '@testing-library/react';

describe('Loader component', () => {
    it('renders loading text', () => {
        const { getByLabelText } = render(<Loader />);
        expect(getByLabelText('Loading...')).toBeInTheDocument();
    });

    it('displays the default loading text', () => {
        const { getByText } = render(<Loader />);
        expect(getByText(LOADING_TEXT)).toBeInTheDocument();
    });
});
