import PropTypes from 'prop-types';
import { getCoverageStatement } from '../../aux/helpers';
import Error from '../Error/Error';
import Loader from '../Loader/Loader';

const List = ({ links, error, loading }) => {
    return (
        <>
            {loading && <Loader />}
            <div role="alert">
                {error && <Error message={error} />}
            </div>

            <div className="list-group">
                <div className="list-group-item list-group-item-action flex-column border-0"></div>
                {links?.map((link, idx) => (
                    <div key={idx} className="list-group-item list-group-item-action flex-column border-0">
                        <div className="row">
                            <h3>
                                <a href={link.target_url} target="_blank" rel="noopener noreferrer">
                                    {link.target_public_name}
                                </a>
                            </h3>
                            <p>{getCoverageStatement(link)}</p>
                        </div>
                    </div>
                ))}
                {(links?.length === 0 || error) && (
                    <>
                        <div className="list-group-item list-group-item-action flex-column border-0">
                            <p>No results found</p>
                        </div>
                    </>
                )}
            </div>
        </>
    );
};

List.propTypes = {
    links: PropTypes.arrayOf(
        PropTypes.shape({
            target_url: PropTypes.string.isRequired,
            target_public_name: PropTypes.string.isRequired
        })
    ),
    error: PropTypes.string,
    loading: PropTypes.bool
};

export default List;
