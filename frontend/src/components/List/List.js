import { useState, useEffect } from 'react';

const List = () => {
    const [elements, setElements] = useState([]);

    const fetchElements= () => {
        fetch('the-new-yorker.json', {
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            }
        })
            .then(response => {
                console.log(response);
                return response.json();
            })
            .then(data => {
                let arrOfLinks = data.ctx_obj[0].ctx_obj_targets[0].target;
                setElements(arrOfLinks);
            })
            .catch(error => {
                console.log(error);
            });
    }
    useEffect(() => {
        fetchElements();
    }, []);

    return (
        <ul>
        {elements && elements.map((element, idx) => <li key={idx}>
            {element.target_public_name}: 
                <a href={element.target_url} target="_blank" rel="noopener noreferrer">{element.target_url.substring(0,70)}...</a>
            </li>)}
        </ul>
    );
};

export default List;
