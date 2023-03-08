const Loader = () => (
    <div className="loading-container">
        {/* TODO: we could put a spinner here: */}
        <div aria-label="Loading...">{LOADING_TEXT}</div>
    </div>
);

export const LOADING_TEXT = 'Loading...';

export default Loader;
