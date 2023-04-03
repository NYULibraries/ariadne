import PropTypes from 'prop-types';

const List = ({ found, links, loading }) => {
    const emptyStyle = { "height": "calc(100vh - 250px)", "width": "100%" };
    return (
        loading ? (
            <div className="empty" style={emptyStyle}></div>
        ) :
            <div className="list-group">
                <div className="list-group-item list-group-item-action flex-column border-0">
                  {found ? "Resource Available Through NYU Libraries" : "Resource Not Available Through NYU Libraries"}
                </div>
                {links?.map((link, idx) => (
                    <div key={idx} className="list-group-item list-group-item-action flex-column border-0">
                        <div className="row">
                            <h3>
                                <a href={link.url} target="_blank" rel="noopener noreferrer">
                                    {link.display_name}
                                </a>
                            </h3>
                            <p>{link.coverage_text}</p>
                        </div>
                    </div>
                ))}
            </div>
    );
};

List.propTypes = {
    found: PropTypes.bool,
    links: PropTypes.arrayOf(
        PropTypes.shape({
            display_name: PropTypes.string.isRequired,
            url: PropTypes.string.isRequired,
        })
    ),
    loading: PropTypes.bool
};

export default List;
