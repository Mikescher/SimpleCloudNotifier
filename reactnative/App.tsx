import { FlashList } from '@shopify/flash-list';
import { StatusBar } from 'expo-status-bar';
import { useRef, useState } from 'react';
import { LayoutAnimation, Pressable, StyleSheet, Text, View } from 'react-native';

export default function App() {
  const generateArray = (size: number) => {
    const arr = new Array(size);
    for (let i = 0; i < size; i++) {
      arr[i] = i;
    }
    return arr;
  };
  
  const [refreshing, setRefreshing] = useState(false);
  const [data, setData] = useState(generateArray(100));
  
  const list = useRef<FlashList<number> | null>(null);
  
  const removeItem = (item: number) => {
    setData(
      data.filter((dataItem) => {
        return dataItem !== item;
      })
    );
    list.current?.prepareForLayoutAnimationRender();
    // after removing the item, we start animation
    LayoutAnimation.configureNext(LayoutAnimation.Presets.easeInEaseOut);
  };
  
  const renderItem = ({ item }: { item: number }) => {
    const backgroundColor = item % 2 === 0 ? "#00a1f1" : "#ffbb00";
    return (
      <Pressable
        onPress={() => {
          removeItem(item);
        }}
      >
        <View
          style={{
            ...styles.container,
            backgroundColor: item > 97 ? "red" : backgroundColor,
            height: item % 2 === 0 ? 100 : 200,
          }}
        >
          <Text>Cell Id: {item}</Text>
        </View>
      </Pressable>
    );
  };
  
  return (
    <FlashList
          ref={list}
          style={{flexGrow: 1, width: '100%', height: 50, backgroundColor: '#88F'}}
          refreshing={refreshing}
          onRefresh={() => {
            setRefreshing(true);
            setTimeout(() => {
              setRefreshing(false);
            }, 2000);
          }}
          keyExtractor={(item: number) => {
            return item.toString();
          }}
          getItemType={(item: number) => {
            return item > 97 ? "even" : "odd";
          }}
          renderItem={renderItem}
          estimatedItemSize={100}
          data={data}
        />
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
