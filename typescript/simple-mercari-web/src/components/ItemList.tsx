import { useEffect, useState } from 'react';
import { Item, fetchItems } from '~/api';

const PLACEHOLDER_IMAGE = import.meta.env.VITE_FRONTEND_URL + '/logo192.png';

// Get image URL
const getImageURL = (imageName: string) => {
  // if there is no image, return placeholder image
  if (!imageName) {
    return PLACEHOLDER_IMAGE;
  }
  return import.meta.env.VITE_BACKEND_URL + '/images/' + imageName;
}


interface Prop {
  reload: boolean;
  onLoadCompleted: () => void;
}

export const ItemList = ({ reload, onLoadCompleted }: Prop) => {
  const [items, setItems] = useState<Item[]>([]);
  useEffect(() => {
    const fetchData = () => {
      fetchItems()
        .then((data) => {
          console.debug('GET success:', data);
          setItems(data.items);
          onLoadCompleted();
        })
        .catch((error) => {
          console.error('GET error:', error);
        });
    };

    if (reload) {
      fetchData();
    }
  }, [reload, onLoadCompleted]);

  return (
    <div className="ItemGrid">
      {items.map((item) => {
        return (
          <div key={item.id} className="ItemList">
            {/* TODO: Task 2: Show item images */}
            <img src={getImageURL(item.image_name)}
            className="Image" />
            <p>
              <span>Name: {item.name}</span>
              <br />
              <span>Category: {item.category}</span>
            </p>
          </div>
        );
      })}
    </div>
  );
};
