import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

# "acousticness"
# "danceability"
# "energy"
# "instrumentalness"
# "key"
# "liveness"
# "loudness"
# "mode"
# "speechiness"
# "tempo"
# "time_signature"
# "valence"

df = pd.read_json("...")

sns.displot(data=df, x="danceability", kind="hist")
plt.show()
