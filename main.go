package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
	GraphQLEndpoint = "https://bff-page.kakao.com/graphql"

	// Thay Cookie mới của bạn vào đây
	CookieString = "_kpdid=fe3277a6e854465f872802cdf1e5891c; _kpiid=278f88bc29a85822679573a2aa34ec0f; _kahai=ba4be1e91c915262b60a6402f5e1e321446b1a1739f6a79ad21e84f1393c8927; _kau=35c5e49a69a9372b6ef6c2bc607965bf7ee6f92bf89bcc5318f18ea5e7dc7503b25394c9f33386b9c72fdd760f46731e48b12b253c105e3355f11c76a1b578a84ae02192398a0a77224fe6fc96b859b3e148c8a8dfb13c5dae2806fe5f7e51ed8478bf2de0240a01b6422618b77f387e6b3a08103120e11a21a054b33336303934353339343736373131303635343630373135333038313634363632e3a084a39eef024d31862ae3336ba294; _kawlt=KoJqnywI4FsAEHmoTpG9DLjSB9i0RC67nFInM6tZNFnium_rbz7Qm6k9wxXXgLmOdVjPOdsgThZ_K2iJtY0nJAbTsrz8ebWiTndnb7RBTZcOt2PtUYbkcIw3o8H_Kyx-; _kawltea=1766718049; _karmt=UsUxsvx5meGQq_TKLAlTO8p1ahRK_2Mjex-jM2BvpqzIA21sR2qKj7wDaWWzSHPr; _karmtea=1769245249; _kaslt=gnpr3z6X3w1YPbfcd1tZmPbLXKQDqzTLwKvruWsvTFrryq0wC921X7LrMZ6wXN3HC8hvN6oVvNhOSrK/qLChRA==; _kpawbat_e=zmUdyUSAugu8jXGNCRoXSmQlO2FDxP%2BVSZgx8QCPkDUTVoVSIRpceNkvQRqQHnCe4ZFVECfKBI9%2BTbWMaDmtKoSjIvCLbQ76hL4QAtlnzDc%3D; _kpwtkn=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..xKgsyFjRhi9D-MJs.KtiWtWIbyG9U0QD4qez5RtdN_u5mrWfrWTk6LdKxEBsdiFfMOG56DiBKMAxWFU66Qx3Sj92A32vJICG87wJt9HkBSi3xSlqxQFuwK3xWZFPY-UOL0GixZkCcWqjaKrVCAqzJfoQxaTV2x19oRsi0nmHTjVEftW6grcPggGwmazD7b5jL2ztoJZ7ZskHfLqmjhEFzGcvhrP3-IKt2z5elrlo__K72nT2zIPzlu7u6AePvYtKge1U6-J4rybFDWlfrHWOs4eEzJy4toK2u5EY7cRXERvWDFDnBxa8ULhRV0qMKRIs7TGjlpO-pzYsD6JLq9g_3hNi0QUHGUch1NPWubc9FBLeza9pfVh6SIZz7kgqlftIymHPgGQOPcJCr_9Xgpf85RGmvuj7Ayx2P3FJFW598dVcASW994mDOVSVZ1V1ygPACiICOH5p2znAs4ChOs0KWB4yJnk4nP5j7HzHtkhbJKPDKM4JBNE1hgdv9hYDb3iMPgwT_GMa4Mt2uCSITvqrvF1i-9W8m4W1HcBnyOZ2Ft8lgyZ1U-fzp4I-ArKa1GeCX7pR7uj1g8MdehsQddDfxjTBusOi13oWdsgeXV3P7odPTJigWM4rzR9XvdC2nPZkaPKk-CYt2gd1Sgni7EMdrOTTbwJLQZX5PPU1_M2cuoHFNSorVhzGUl1wR4z2RLmveGFHG-7vCKk06qbZunO66-6PBUkwEYOJrtcjA54nGzOMJpevDSXNcMVjLxPDIf-a056mqq27s-_2v4XRkIHC5LbXIkC-wje-eAf7gue9fgCAQ8kXDdoggiicCWF968A3SrW8aQpKXJIkcytJN5xbBFECgHnT7PeLVAKFPYyuH8KCcpLK2nal4cOhXYBtPSlNd_5ybx7ww6IVGdG32iSFYTyG8QTEHRm8mne0BVvQqodivUJzUbkUgFLxbvATWiQHlK4cCbGgLNflrDV_HXXZAmake0ZI0LT8dZlPkAK19poW-qiU6oJ4Oge-T6H8nSqDZI0EY_-4YXGH1gMPNxImGJyUdzwLK809q49RxyUx2yBkMhMtQzu5J7iARLqorXKZPt11tpryxsoNfn1ROZ8al5dgs8yGGuW6cXyf2LTlu17p0pigj7OJuqcL8H3DlyjbHMuGQQu3SqyZSa54YzCu9iZNt8se9b3Z1_bnZvn1CrhSgnNJ4vp0nt9hkd4DGBFPCI-g9xockSodiJEUqDBjLW7_9uj9hrn-iaKwYxWoAEHbUfkpqz9a54FP7yAxYIIVxRO3ck5HNgONKbHfPxLr1-8BFXrkfZ8S80cB5UkN-gCzm4BMOIwhuZ352AGtqqhZZ4ertvY40d-cvs35zGijzDYm1bFxE8ntwlY5TQKK5GtarU_Ml2NiA76yKzKahURApYcGuYixwqbOuGNjHuc5qQlLct0ECJe3l6oT7InbuGW0-BrCjbb4kbiRk9aa2XiEayNNLdnkiPNsQc4p1DzmilSiav6Gdo7NDqjY1FI00T8Rm367z7--CVwSy2klXRFaAslo1eDO-6oTVmhMvcsgkCe5I4IxYoIqXSFnZ2XiLGhcBp8s9xGcb2RhwcbfM0iCHELyQzpufOoKkrd1_ZwPh_LtYHN2r8_FJjTcf0pgRzymoqL31y5HMwzgg7ss7innQRDi4vKeprtuiscB9PqZcbCJY3_gF2aBc1kaYIl5EsrLx2oqDpmhKjhCULmqFXauJpa2toYAsxpkdOcGJZOUeoiVMf0KhPHlaVNJHwvh_-97uoy5_GKcdUBQ0bDI02m7MQfpqbvxQLZAQ81KGnBn5rzcBKnPrCEVKtN3a2firjtYzcG8hIsoVxX868c-IuYW4IWn4h3yQ5Dj-BcFPtZbfv0hsdjFvX4FBWRktOneINVAvcIhj0EHP9GeoGewITpGjYcd-9uUXAAovDQedZhv0Ive0QPwRrs2ENnDIgRW5YesuyLXD1-u-wYC5f-zh0LcwjL9u8wIwgPjxTyYSp6CbQnZm6S967ceErKV5gr8e4-Ac7PGS0Am7CqQ4DVANnsd1cd3wTcV2b1MdcrIy6BOhbCU80KpRD9uNratNKLSFXloNe5VO2RNFtJgxncWlxH8872Gmbu--GQOWeD1qsoy6Av5j8wvhjrqTa4jGCbx4urlX_1KtJ_SLBavllc0HYUgabTK5j5zPMFKEuhunYgANbsr1zC8OABZ-q3c_RPRNJpqx9JqjF3RNPWWGRQzshQw63SPEec-iB_CKY2_hHTWWZxje3hoK55F6U-cApDX2AfZYyTMpW32bA6irlDtN44UEQfsV2swtIAkpk4zkthUy8Ys.Cvrpz_Tpgweo-TyLh6ithw; _T_ANO=i7xIOkQp8fWV8kJZr4TdwwgqQtEvgD/aDxT8IwpnIogr7u9bgnMb3RygMwX6x84KDghfkPKRCAXWWwcw3cMkg84hSz2JZ8tcYVanXuLhqJ/yfb/HyJEuTFNGaJ/r590K1Sk0v/DgadRRrRO5gvKNCRYsG4E4M/fNoclA0sOmt6iB1Kzn3T+KR6B6rUPvTaacL6DW4sNK46XLRIbCkR+51ytW+7tXgzXabepMdydcF9jMLbOcJPA+/2IZMJzQQ0S/teIbyczn7O/Cp+PzaOnv6bzJKVY99SNgrrY5cHMyuOIUEBGE95caJT2xiXu+aQFZR0smGeYz7FzVn1K18AxGYA=="
)

// ============================================================================
// QUERY (dùng payload bạn cung cấp — giữ nguyên các type Long!)
// ============================================================================

const ViewerQuery = `
query viewerInfo($seriesId: Long!, $productId: Long!) {
  viewerInfo(seriesId: $seriesId, productId: $productId) {
    item {
      ...SingleFragment
    }
    seriesItem {
      ...SeriesFragment
    }
    prevItem {
      ...NearItemFragment
    }
    nextItem {
      ...NearItemFragment
    }
    viewerData {
      ...TextViewerData
      ...TalkViewerData
      ...ImageViewerData
      ...VodViewerData
    }
    displayAd {
      ...DisplayAd
    }
  }
}

fragment SingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  series {
    ...SeriesFragment
  }
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  discountRate
  discountRateText
  isShortsDrama
}

fragment SeriesFragment on Series {
  id
  seriesId
  title
  thumbnail
  landThumbnail
  categoryUid
  lang
  category
  categoryType
  subcategoryUid
  subcategory
  badge
  isAllFree
  isWaitfree
  ageGrade
  state
  onIssue
  authors
  description
  pubPeriod
  freeSlideCount
  lastSlideAddedDate
  waitfreeBlockCount
  waitfreePeriodByMinute
  bm
  saleState
  startSaleDt
  saleMethod
  discountRate
  discountRateText
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  translateProperty {
    ...TranslatePropertyFragment
  }
}

fragment ServicePropertyFragment on ServiceProperty {
  viewCount
  readCount
  ratingCount
  ratingSum
  commentCount
  pageContinue {
    ...ContinueInfoFragment
  }
  todayGift {
    ...TodayGift
  }
  preview {
    ...PreviewFragment
  }
  waitfreeTicket {
    ...WaitfreeTicketFragment
  }
  isAlarmOn
  isLikeOn
  ticketCount
  purchasedDate
  lastViewInfo {
    ...LastViewInfoFragment
  }
  purchaseInfo {
    ...PurchaseInfoFragment
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment ContinueInfoFragment on ContinueInfo {
  title
  isFree
  productId
  lastReadProductId
  scheme
  continueProductType
  hasNewSingle
  hasUnreadSingle
}

fragment TodayGift on TodayGift {
  id
  uid
  ticketType
  ticketKind
  ticketCount
  ticketExpireAt
  ticketExpiredText
  isReceived
  seriesId
}

fragment PreviewFragment on Preview {
  item {
    ...PreviewSingleFragment
  }
  nextItem {
    ...PreviewSingleFragment
  }
  usingScroll
}

fragment PreviewSingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
}

fragment OperatorPropertyFragment on OperatorProperty {
  thumbnail
  copy
  helixImpId
  isTextViewer
  selfCensorship
  isBook
  cashInfo {
    discountRate
    setDiscountRate
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment AssetPropertyFragment on AssetProperty {
  bannerImage
  cardImage
  cardTextImage
  cleanImage
  ipxVideo
  bannerSet {
    ...BannerSetFragment
  }
  cardSet {
    ...CardSetFragment
  }
  cardCover {
    ...CardCoverFragment
  }
}

fragment BannerSetFragment on BannerSet {
  backgroundImage
  backgroundColor
  mainImage
  titleImage
}

fragment CardSetFragment on CardSet {
  backgroundColor
  backgroundImage
}

fragment CardCoverFragment on CardCover {
  coverImg
  coverRestricted
}

fragment WaitfreeTicketFragment on WaitfreeTicket {
  chargedPeriod
  chargedCount
  chargedAt
}

fragment LastViewInfoFragment on LastViewInfo {
  isDone
  lastViewDate
  rate
  spineIndex
}

fragment PurchaseInfoFragment on PurchaseInfo {
  purchaseType
  rentExpireDate
  expired
}

fragment TranslatePropertyFragment on TranslateProperty {
  category {
    ...LocaleMapFragment
  }
  sub_category {
    ...LocaleMapFragment
  }
}

fragment LocaleMapFragment on LocaleMap {
  ko
  en
  th
}

fragment NearItemFragment on NearItem {
  productId
  slideType
  ageGrade
  isFree
  title
  thumbnail
}

fragment TextViewerData on TextViewerData {
  type
  atsServerUrl
  metaSecureUrl
  contentsList {
    chapterId
    contentId
    secureUrl
  }
}

fragment TalkViewerData on TalkViewerData {
  type
  talkDownloadData {
    dec
    host
    path
    talkViewerType
  }
}

fragment ImageViewerData on ImageViewerData {
  type
  imageDownloadData {
    ...ImageDownloadData
  }
}

fragment ImageDownloadData on ImageDownloadData {
  files {
    ...ImageDownloadFile
  }
  totalCount
  totalSize
  viewDirection
  gapBetweenImages
  readType
}

fragment ImageDownloadFile on ImageDownloadFile {
  no
  size
  secureUrl
  width
  height
}

fragment VodViewerData on VodViewerData {
  type
  vodDownloadData {
    contentPackId
    drmType
    endpointUrl
    width
    height
    duration
  }
  drmInfo {
    type
    serverType
    error
    fairplayCertificateUrl
    widevineLicenseUrl
    fairplayLicenseUrl
    token
    provider
    assertion
  }
}

fragment DisplayAd on DisplayAd {
  sectionUid
  bannerUid
  treviUid
  momentUid
}
`

// ============================================================================
// MODELS (mình chỉ khai báo những field cần dùng)
// ============================================================================

type GraphQLPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type GQLResponse struct {
	Data struct {
		ViewerInfo ViewerInfo `json:"viewerInfo"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ViewerInfo struct {
	Item       Item       `json:"item"`
	SeriesItem *Series    `json:"seriesItem,omitempty"`
	PrevItem   *NearItem  `json:"prevItem,omitempty"`
	NextItem   *NearItem  `json:"nextItem,omitempty"`
	ViewerData ViewerData `json:"viewerData"`
}

type Item struct {
	Title string `json:"title"`
	// thêm field nếu cần về sau
	ProductId int `json:"productId,omitempty"`
	SeriesId  int `json:"seriesId,omitempty"`
}

type Series struct {
	SeriesId int    `json:"seriesId"`
	Title    string `json:"title"`
}

type NearItem struct {
	ProductId int    `json:"productId"`
	Title     string `json:"title"`
	// slideType, etc. omitted — thêm nếu cần
}

type ViewerData struct {
	AtsServerUrl string        `json:"atsServerUrl"`
	MetaSecure   string        `json:"metaSecureUrl,omitempty"`
	ContentsList []ContentPart `json:"contentsList"`
}

type ContentPart struct {
	ChapterId int    `json:"chapterId,omitempty"`
	ContentId int    `json:"contentId,omitempty"`
	SecureUrl string `json:"secureUrl"`
}

type ContentFileResponse struct {
	ContentInfo struct {
		ParagraphList []ParagraphNode `json:"paragraphList"`
	} `json:"contentInfo"`
}

type ParagraphNode struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Children []ParagraphNode `json:"childParagraphList,omitempty"`
}

// ============================================================================
// MAIN: loop theo nextItem.productId — giữ flow tải/chạy y nguyên code gốc
// ============================================================================

func main() {
	// context tổng cho toàn bộ chạy
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	headers := map[string]string{
		"accept":          "*/*",
		"accept-language": "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7",
		"content-type":    "application/json",
		"origin":          "https://page.kakao.com",
		"referer":         "https://page.kakao.com/",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
		"cookie":          CookieString,
	}

	// Bắt đầu từ đây: set seriesId và productId chương ban đầu bạn muốn crawl
	currentSeries := 61097897
	currentProduct := 66447436

	for {
		fmt.Printf("\n>>> Crawl productId = %d (seriesId=%d)\n", currentProduct, currentSeries)

		// 1) build payload giữ nguyên như bạn yêu cầu
		payload := GraphQLPayload{
			OperationName: "viewerInfo",
			Query:         ViewerQuery,
			Variables: map[string]interface{}{
				"seriesId":  currentSeries,
				"productId": currentProduct,
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		respBody, err := makeRequest(ctx, client, "POST", GraphQLEndpoint, headers, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf(">>> ERROR: GraphQL request failed: %v\n", err)
			break
		}

		var gqlResp GQLResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			fmt.Printf(">>> ERROR: Parse Error: %v\nBody: %s\n", err, string(respBody))
			break
		}
		if len(gqlResp.Errors) > 0 {
			fmt.Printf(">>> ERROR: GraphQL Error: %s\n", gqlResp.Errors[0].Message)
			break
		}

		viewerData := gqlResp.Data.ViewerInfo.ViewerData
		itemInfo := gqlResp.Data.ViewerInfo.Item
		nextItem := gqlResp.Data.ViewerInfo.NextItem

		if viewerData.AtsServerUrl == "" {
			fmt.Println(">>> ERROR: Không tìm thấy atsServerUrl. Kiểm tra lại Cookie.")
			break
		}

		// Tạo file
		safeTitle := sanitizeFilename(itemInfo.Title)
		fileName := safeTitle + ".txt"

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf(">>> ERROR: Không thể tạo file %s: %v\n", fileName, err)
			break
		}
		fmt.Printf(">>> Đã tạo file: %s\n", fileName)
		fmt.Printf(">>> Tổng số phần: %d. Bắt đầu tải...\n", len(viewerData.ContentsList))

		// Ghi tiêu đề
		if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", itemInfo.Title)); err != nil {
			fmt.Printf(">>> WARNING: lỗi ghi tiêu đề: %v\n", err)
		}

		// Loop tải các phần
		for i, part := range viewerData.ContentsList {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}
			fullUrl := viewerData.AtsServerUrl + part.SecureUrl
			fmt.Printf(" -> [%d/%d] Đang tải... (Delay 1s)\n", i+1, len(viewerData.ContentsList))

			partBody, err := makeRequest(ctx, client, "GET", fullUrl, headers, nil)
			if err != nil {
				fmt.Printf("   [Lỗi tải phần %d]: %v\n", i+1, err)
				continue
			}

			var contentResp ContentFileResponse
			if err := json.Unmarshal(partBody, &contentResp); err != nil {
				// nếu không parse được JSON -> có thể là image / binary -> bỏ qua
				continue
			}

			for _, node := range contentResp.ContentInfo.ParagraphList {
				text := extractTextRecursive(node)
				text = strings.ReplaceAll(text, "&nbsp;", " ")
				if _, err := file.WriteString(text); err != nil {
					fmt.Printf("   Lỗi ghi file: %v\n", err)
				}
			}
		}

		file.Close()
		fmt.Println(">>> XONG! Kiểm tra file:", fileName)

		// Nếu không có nextItem -> dừng
		if nextItem == nil || nextItem.ProductId == 0 {
			fmt.Println(">>> HẾT CHƯƠNG hoặc không tìm thấy nextItem.")
			break
		}

		// Cập nhật productId để crawl chương tiếp theo
		fmt.Printf(">>> Next: %d (%s)\n", nextItem.ProductId, nextItem.Title)
		currentProduct = nextItem.ProductId

		// delay an toàn trước lần tiếp theo
		time.Sleep(2 * time.Second)
	}
}

// ============================================================================
// HELPERS
// ============================================================================

func makeRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status %d | Body: %s", resp.StatusCode, string(respBytes))
	}
	return respBytes, nil
}

func extractTextRecursive(node ParagraphNode) string {
	var sb strings.Builder
	if node.Text != "" {
		sb.WriteString(node.Text)
	}
	if node.Type == "BR" {
		sb.WriteString("\n")
	}
	for _, child := range node.Children {
		sb.WriteString(extractTextRecursive(child))
	}
	switch node.Type {
	case "P", "DIV":
		sb.WriteString("\n")
	case "H3":
		sb.WriteString("\n\n### ")
	}
	return sb.String()
}

func sanitizeFilename(name string) string {
	reg, _ := regexp.Compile(`[\\/:*?"<>|]`)
	safe := reg.ReplaceAllString(name, "_")
	return strings.TrimSpace(safe)
}
