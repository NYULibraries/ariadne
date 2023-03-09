import PropTypes from 'prop-types';
import { getCoverageStatement } from '../../aux/helpers';

const List = ({ links, loading }) => {
    const emptyStyle = { "height": "calc(100vh - 250px)", "width": "100%" };
    return (
        loading ? (
            <div className="empty" style={emptyStyle}></div>
        ) :
            <div className="list-group">
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
            </div>
    );
};

List.propTypes = {
    links: PropTypes.arrayOf(
        PropTypes.shape({
            target_url: PropTypes.string.isRequired,
            target_public_name: PropTypes.string.isRequired
        })
    ),
    loading: PropTypes.bool
};

export default List;
